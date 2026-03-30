# Build fat-jar (Maven Shade) i lekki runtime JRE.
FROM maven:3.9.9-eclipse-temurin-21-alpine AS build
WORKDIR /app
COPY pom.xml .
RUN mvn -q -B dependency:go-offline
COPY src ./src
RUN mvn -q -B package -DskipTests

FROM eclipse-temurin:21-jre-alpine
WORKDIR /app
COPY --from=build /app/target/java-demo-worker.jar /app/worker.jar
USER 65534:65534
ENTRYPOINT ["java", "-jar", "/app/worker.jar"]
