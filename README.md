# Buttoners

## Схема сервисов

![C4](./scheme.drawio.png)

## Схема инфраструктуры

![C4](./scheme_infr.drawio.png)

mermaid версия

```mermaid
C4Context
%% Для "красивого" рендера пришлось изменить + добавить пвсевдо границы
UpdateLayoutConfig($c4ShapeInRow="3", $c4BoundaryInRow="1")

Person(admin, "Администратор", "Главный защитник кнопки")

System_Ext(grafana, "Grafana", "Визуализация метрик, логов, трейсов")

Boundary(b1, "", ""){
    System_Ext(metrics, "Victoria Metrics", "Хранение метрик")
    System_Ext(logs, "Victoria Logs", "Хранение логов")
    System_Ext(traces, "Victoria Traces", "Хранение трейсов")
}

Boundary(b2, "", ""){
    System_Ext(collector, "Grafana Alloy", "Хранение трейсов")
    System_Ext(docker, "Docker", "Система контейнеризации")
}

Boundary(b3, "", ""){
    System(service, "[Container: Go]", "Любой кнопочный сервис")
}

Rel(admin, grafana, "Смотрит метрики, логи, трейсы")

Rel(grafana, metrics, "Получает метрики для просмотра", "JSON/HTTP")
Rel(grafana, logs, "Получает логи для просмотра", "JSON/HTTP")
Rel(grafana, traces, "Получает трейсы для просмотра", "JSON/HTTP")

Rel(collector, service, "Собирает метрики", "JSON/HTTP")
Rel(collector, docker, "Собирает логи", "unix")
Rel(service, collector, "Отправляет трейсы", "gRPC")


Rel(collector, metrics, "Отправляет метрики", "JSON/HTTP")
Rel(collector, logs, "Отправляет логи", "JSON/HTTP")
Rel(collector, traces, "Отправляет трейсы", "JSON/HTTP")
```

## Запустил и открыл по быстрому

- [Метрики](http://metrics.localhost)
- [Логи](http://logs.localhost)
- [Трейсы](http://traces.localhost)
- [Grafana](http://grafana.localhost)
- [Legacy](http://legacy.localhost)
- [Traefik](http://traefik.localhost)
- [Данные в БД](http://db.localhost)
- [Alloy](http://alloy.localhost)
