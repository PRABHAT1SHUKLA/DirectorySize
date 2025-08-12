receivers:
  hostmetrics:
    collection_interval: 10s
    scrapers:
      cpu:
      memory:
      disk:
      filesystem:
      network:

exporters:
  logging:
    loglevel: debug

service:
  pipelines:
    metrics:
      receivers: [hostmetrics]
      exporters: [logging]
