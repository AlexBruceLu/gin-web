# gin-web

## 链路追踪

1. **为啥要用链路追踪?**

   > 在微服务架构中,一个请求,请求了多个服务单元,如果请求出现了错误,很难定位错误,所以这时候就需要链路追踪,来定位bug

2. **使用三方库**

   - Jaeger：https://www.jaegertracing.io
   - Zipkin：https://zipkin.io/
   - Appdash：https://about.sourcegraph.com/

   > 以jaeger为例,linux安装
   >
   > ```bash
   > $ wget -c https://github.com/jaegertracing/jaeger/releases/download/v1.15.1/jaeger-1.15.1-linux-amd64.tar.gz
   > ```

