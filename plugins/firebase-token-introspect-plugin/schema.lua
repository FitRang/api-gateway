return {
	name = "firebase-token-introspect-plugin",
	fields = {
		{
			config = {
				type = "record",
				fields = {
					{
						firebase_project_id = {
							type = "string",
							required = true,
						},
					},
					{
						header_user_id = {
							type = "string",
							required = false,
							default = "X-User-ID",
						},
					},
					{
						header_user_email = {
							type = "string",
							required = false,
							default = "X-User-Email",
						},
					},
				},
			},
		},
	},
}
