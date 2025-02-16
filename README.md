# Background Worker API

Worker que realiza llamadas HTTP periódicas a endpoints configurables y expone los logs vía API.

## Configuración

1. Crea un archivo `.env` con la configuración de las APIs:
```json
APIS=[
    {
        "url": "https://tu-api.com/ping",
        "interval": "5s"
    }
]
```

Formatos de intervalo soportados:
- "5s" (5 segundos)
- "1m" (1 minuto)
- "1h" (1 hora)

## Instalación

```bash
# Clonar repositorio
git clone <repositorio>

# Instalar dependencias
go mod tidy
```

## Ejecución

```bash
go run app.go
```

## Endpoints

### GET /logs
Retorna los logs del worker. Requiere autenticación.

Headers requeridos:
- X-API-Key: Tu clave API definida en .env

Ejemplo:
```bash
curl -H "X-API-Key: tu-clave-secreta-aqui" http://localhost:8080/logs
```

## Estructura del Proyecto

```
.
├── api.log          # Archivo de logs
├── app.go           # Punto de entrada
├── internal/        # Lógica interna
│   ├── config.go    # Configuración
│   ├── handler.go   # Handlers HTTP 
│   └── worker.go    # Lógica del worker
└── README.md
```

## Logs

Los logs usan emojis para indicar el estado:
- ✅ Éxito
- ❌ Error
- ⚠️ Advertencia