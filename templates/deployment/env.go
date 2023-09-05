package deployment

func EnvFileBytes() []byte {
	return []byte(
		`SUPABASE_URL=
SUPABASE_KEY=
SURREALDB_URL=ws://localhost:8000/rpc`)
}
