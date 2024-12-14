### Настройка переменных окружения

.env.example -> .env

### Запуск сервера

```bash
go mod download
go run exporter/cmd/exporter/main.go
```

### Запросы PromQL

```
{__name__=~"metric_(cpu_usage|memory_usage|disk_usage)"}
```

```
metric_cpu_usage
```

```
metric_memory_usage
```

```
metric_memory_usage
```
