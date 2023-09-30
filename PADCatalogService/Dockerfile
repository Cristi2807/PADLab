# Use the official .NET SDK image as the base image
FROM mcr.microsoft.com/dotnet/sdk:6.0 AS build

# Set the working directory inside the container
WORKDIR /app/catalog

COPY catalog/catalog.fsproj .
RUN dotnet restore
COPY catalog/. .
RUN dotnet build -c Release


WORKDIR /app/catalog.dal

COPY catalog.dal/catalog.dal.fsproj .
RUN dotnet restore
COPY catalog.dal/. .
RUN dotnet build -c Release


# Copy the .NET project file(s) and restore dependencies
WORKDIR /app/catalog.rest
COPY catalog.rest/catalog.rest.fsproj .
RUN dotnet restore

# Copy the rest of the application code
COPY catalog.rest/. .

# Build the application
RUN dotnet publish -c Release -o out

# Use the official ASP.NET Core runtime image as the base image for the final image
FROM mcr.microsoft.com/dotnet/aspnet:6.0 AS runtime

# Set the working directory inside the container
WORKDIR /app

# Copy the published output from the build image to the runtime image
COPY --from=build /app/catalog.rest/out .

# Expose the port your Giraffe app is listening on (change 8080 to the appropriate port)
EXPOSE 5050

# Start the application
ENTRYPOINT ["dotnet", "catalog.rest.dll"]
