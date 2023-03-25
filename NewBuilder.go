package sql

func NewBuilder(dialect string) *Builder {
	builderSql := map[string]any{
		"columns": []map[string]any{},
	}
	return &Builder{
		Dialect: dialect,
		sql:     builderSql,
	}
}
