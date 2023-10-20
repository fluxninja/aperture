# syntax=docker/dockerfile:1.4

FROM mcr.microsoft.com/dotnet/sdk:7.0 AS builder

WORKDIR /app

COPY ApertureSDK.sln .
COPY ApertureSDK.csproj .
COPY Examples/Examples.csproj ./Examples/

RUN dotnet restore

COPY . .

RUN dotnet build -c Release -o /app/build

FROM builder AS publish
RUN dotnet publish -c Release -o /app/publish Examples/Examples.csproj

FROM mcr.microsoft.com/dotnet/aspnet:7.0 AS final
WORKDIR /app

RUN apt-get update \
  && apt-get install -y --no-install-recommends \
  ca-certificates \
  wget \
  && apt-get clean

COPY --from=publish /app/publish .
COPY Examples/log4net.config .

HEALTHCHECK --interval=5s --timeout=60s --retries=3 --start-period=5s \
  CMD wget --no-verbose --tries=1 --spider 127.0.0.1:8080/health || exit 1

ENTRYPOINT ["dotnet", "Examples.dll"]
