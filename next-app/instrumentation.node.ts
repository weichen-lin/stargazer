import { NodeSDK } from '@opentelemetry/sdk-node'
import { OTLPTraceExporter } from '@opentelemetry/exporter-trace-otlp-http'
import { Resource } from '@opentelemetry/resources'
import { SEMRESATTRS_SERVICE_NAME } from '@opentelemetry/semantic-conventions'
import { SimpleSpanProcessor } from '@opentelemetry/sdk-trace-node'

const sdk = new NodeSDK({
  resource: new Resource({
    [SEMRESATTRS_SERVICE_NAME]: 'next-app',
  }),
  spanProcessor: new SimpleSpanProcessor(
    new OTLPTraceExporter({
      url: 'https://otel.kloudmate.com:4318/v1/traces',
      headers: {
        Authorization: 'sk_sAUT2yCaEZ1E3nzOgr4wddBz',
      },
    }),
  ),
})

sdk.start()
