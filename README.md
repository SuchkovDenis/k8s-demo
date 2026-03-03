# Kubernetes за 20 минут — код и материалы

## Структура проекта

```
k8s-demo/
├── app/
│   ├── main.go          # Go-сервис с CPU-нагрузкой
│   ├── go.mod
│   └── Dockerfile       # Multi-stage build
├── k8s/
│   ├── deployment.yaml  # Deployment с resource limits
│   ├── service.yaml     # NodePort Service
│   └── hpa.yaml         # HorizontalPodAutoscaler
├── dashboard/
│   ├── index.html       # Веб-дашборд для нагрузочного тестирования
│   └── server.py        # HTTP-сервер с проксированием K8s API
└── README.md
```

## Инструменты

- minikube v1.38.1
- kubectl v1.35.1
- Docker v27.4.0
- hey (HTTP load testing)

## Быстрый старт

### 1. Сборка образа
```bash
cd app
docker build -t demo-app:1.0 .
```

### 2. Демо: Docker с ограничениями
```bash
docker run -d --name demo-app -p 8080:8080 --cpus=0.2 --memory=128m demo-app:1.0

# Открыть дашборд (работает и для Docker-режима)
cd dashboard && python3 server.py
# http://localhost:3000 → ввести http://localhost:8080 → START

docker rm -f demo-app
```

### 3. Демо: Kubernetes с автоскейлингом
```bash
minikube start
minikube addons enable metrics-server
minikube image load demo-app:1.0

kubectl apply -f k8s/deployment.yaml -f k8s/service.yaml -f k8s/hpa.yaml
minikube service demo-app --url  # получить URL

# Запустить kubectl proxy + дашборд
kubectl proxy --port=8001 &
cd dashboard && python3 server.py
# http://localhost:3000 → ввести URL сервиса → START
```

### 4. Наблюдение за автоскейлингом
```bash
kubectl get hpa -w          # следить за HPA
kubectl get pods -w         # следить за подами
kubectl top pods            # CPU/Memory
```
