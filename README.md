# Web Scraper API - TDL

API desarrollada en GO, como trabajo práctico para la materia: Teoría de Lenguaje

## Integrantes

- Martin Alejo Polese - 106808
- Lucas Grati - 102676
- Edgardo Francisco Saez - 104896

## Scripts

Al posicionarte en el directorio correspondiente, se pueden ejecutar los siguientes comandos:

### Frontend

```bash
npm install
```

Instala las dependencias necesarias para levantar el frontend.

```bash
npm start
```

Levanta el frontend en modo desarrollo.\
Abrir [http://localhost:3000](http://localhost:3000) para verlo en el navegador.

### Backend

```bash
go run .
```

Levanta el servidor en el puerto 8080.

## Endpoints

Todas las llamadas a los endpoints son con método ***GET***.\
Recordar de anteponer [http://localhost:8080](http://localhost:8080) en las llamadas a la API.

### General

Scrapea notebooks de Mercado Libre, Fravega y Fullh4rd

```bash
/api/general
```

### Mercado Libre

Scrapea notebooks de Mercado Libre

```bash
/api/mercadolibre
```

### Fravega

Scrapea notebooks de Fravega

```bash
/api/fravega
```

### FullH4rd

Scrapea notebooks de FullH4rd

```bash
/api/fullh4rd
```

## Query parameters

### RAM:

- MinRam
- MaxRam

### Pulgadas:

- MinInches
- MaxInches

### SSD:

- MinStorage
- MaxStorage

### Precio:

- MinPrice
- MaxPrice

### Procesador:

- Processor

    ***Nota**: Puede ser "amd", "intel" o "apple"*