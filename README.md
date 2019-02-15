# metrics-generator

Randomly generates Prometheus metrics for simulating default /metrics endpoints accross servers

## Usage

1. **Build**

```
docker-comopse build
```

2. **Run**

```
docker-compose up -d
```

3. **Test**

```
curl http://localhost:32865/metrics
```

## Cause accidents!

If you want to cause any resource metric to behave abnormaly, you could:

### 1. Cause latency accidents 
```
curl -X POST http://localhost:32865/accidents
{
	"resourcename": "/resource/test-0001",
	"type": "latency",
	"value": "100"
}
```

* Which will cause the requests time counter to the resource `/resource/test-0001` to have a median duration of 100ms

### 2. Cause request count accidents
```
curl -X POST http://localhost:32865/accidents
{
	"resourcename": "/resource/test-0001",
	"type": "calls",
	"value": "100"
}
```

* Which will cause the request count to the resource `/resource/test-0001` to increase by a factor of 100

### 3. Cause error rate accidents
```
curl -X POST http://localhost:32865/accidents
{
	"resourcename": "/resource/test-0001",
	"type": "errorrate",
	"value": "0.7"
}
```

* Which will cause the error rate associated with the resource `/resource/test-0001` to float around 70%


**To remove all accidents, send a HTTP DELETE to the same endpoint**

```
curl -X DELETE http://localhost:32865/surgery-accident
```

## Control /metrics entropy

By default, this generator simulates requests to 10 different URIs, 2 different service versions, 2 different app versions and 2 different device types

Via the `/entropy/set` endpoint you can control this numbers resulting in how many samples are collected each time.

```
curl - X POST http://localhost:32865/entropy/set
{
	"uricount": 20,
	"serviceversioncount": 2,
	"appversioncount": 2,
	"devicecount": 2
}
```

* This request will cause the generator to simulate requests to 20 different URIs
