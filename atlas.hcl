env "local" {
  url = "mysql://root:${getenv("MYSQL_ROOT_PASSWORD")}@localhost:3306/${getenv("MYSQL_DATABASE")}"
  dev = "docker://mysql/8/test"
  src = "ent://internal/infra/ent/schema"
  migration {
    dir = "file://internal/infra/ent/migrate/migrations"
  }
}
