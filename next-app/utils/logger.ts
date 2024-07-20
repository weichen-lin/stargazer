import { createLogger, transports } from 'winston'
import { OTLPLogExporter } from '@opentelemetry/exporter-logs-otlp-http'
import { LoggerProvider, BatchLogRecordProcessor } from '@opentelemetry/sdk-logs'
import { Resource } from '@opentelemetry/resources'
import { SeverityNumber } from '@opentelemetry/api-logs'

const resource = Resource.default().merge(
  new Resource({
    'service.name': 'next-app',
    'service.version': '0.1.0',
  }),
)

const loggerProvider = new LoggerProvider({
  resource: resource,
})
const logExporter = new OTLPLogExporter({
  url: `https://otel.kloudmate.com:4318/v1/logs`,
  headers: {
    Authorization: process.env.KLOUDMATE_TOKEN,
  },
})
const logProcessor = new BatchLogRecordProcessor(logExporter)
loggerProvider.addLogRecordProcessor(logProcessor)

const formatLog = (args: any) => (typeof args === 'string' ? args : JSON.stringify(args))

const consoleTransport = new transports.Console()
const logger = createLogger({
  transports: [consoleTransport],
})

const Logger = {
  ...logger,
  info: (args: any) => {
    loggerProvider.getLogger('next-app').emit({ body: formatLog(args), severityNumber: SeverityNumber.INFO })
    return logger.info(args)
  },

  error: (args: any) => {
    loggerProvider.getLogger('next-app').emit({ body: formatLog(args), severityNumber: SeverityNumber.ERROR })
    return logger.error(args)
  },
}

export default Logger
