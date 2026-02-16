local jwt = require("resty.jwt")
local http = require("resty.http")
local cjson = require("cjson.safe")
local kong = kong
local ngx = ngx

local FirebaseTokenIntrospectPlugin = {
	PRIORITY = 1000,
	VERSION = "1.0.0",
}

local JWKS_CACHE_KEY = "firebase_jwks"
local JWKS_CACHE_TTL = 3600

local function fetch_jwks_from_google()
	local httpc = http.new()

	local res, err =
		httpc:request_uri("https://www.googleapis.com/service_accounts/v1/jwk/securetoken@system.gserviceaccount.com", {
			method = "GET",
			ssl_verify = true,
		})

	if not res then
		return nil, err
	end

	return cjson.decode(res.body)
end

local function get_jwks()
	local jwks, err = kong.cache:get(JWKS_CACHE_KEY, { ttl = JWKS_CACHE_TTL }, fetch_jwks_from_google)

	return jwks, err
end

local function extract_token()
	local auth_header = kong.request.get_header("authorization")

	if auth_header then
		local m, err = ngx.re.match(auth_header, "Bearer\\s+(.+)", "jo")
		if m then
			return m[1]
		end
	end

	local token = kong.request.get_query_arg("token")
	return token
end

local function find_key(jwks, kid)
	if not jwks or not jwks.keys then
		return nil
	end

	for _, key in ipairs(jwks.keys) do
		if key.kid == kid then
			return jwt:jwk_to_pem(key)
		end
	end

	return nil
end

local function verify_jwt(token, project_id)
	local jwt_obj = jwt:load_jwt(token)

	if not jwt_obj.valid then
		return nil, "invalid jwt"
	end

	local kid = jwt_obj.header.kid

	local jwks, err = get_jwks()

	if not jwks then
		return nil, err
	end

	local pem = find_key(jwks, kid)

	if not pem then
		return nil, "public key not found"
	end

	local verified = jwt:verify_jwt_obj(pem, jwt_obj)

	if not verified.verified then
		return nil, "signature verification failed"
	end

	local payload = verified.payload

	if payload.aud ~= project_id then
		return nil, "invalid audience"
	end

	if payload.iss ~= "https://securetoken.google.com/" .. project_id then
		return nil, "invalid issuer"
	end

	if payload.exp < ngx.time() then
		return nil, "token expired"
	end

	return payload
end

function FirebaseTokenIntrospectPlugin:access(conf)
	local token = extract_token()

	if not token then
		return kong.response.exit(401, {
			error = "missing_token",
		})
	end

	local payload, err = verify_jwt(token, conf.firebase_project_id)

	if err then
		return kong.response.exit(401, {
			error = err,
		})
	end

	if payload.email_verified ~= true then
		return kong.response.exit(403, {
			error = "email_not_verified",
		})
	end

	kong.service.request.set_header(conf.header_user_id, payload.user_id)
	kong.service.request.set_header(conf.header_user_email, payload.email)
end
