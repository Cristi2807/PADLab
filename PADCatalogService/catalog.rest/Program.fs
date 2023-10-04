open Giraffe
open Microsoft.AspNetCore.Builder
open Microsoft.AspNetCore.Hosting
open Microsoft.Extensions.DependencyInjection
open Microsoft.Extensions.Hosting
open catalog.rest
open System.Net.Http
open System

let webApp =
  choose [ ProductApi.productRoutes; RequestErrors.NOT_FOUND "Not found" ]

let configureApp (app: IApplicationBuilder) = app.UseGiraffe webApp

let configureServices (services: IServiceCollection) = services.AddGiraffe() |> ignore

[<EntryPoint>]
let main _ =

  let serviceDiscoveryURL =
    match Environment.GetEnvironmentVariable "SERVICE_DISCOVERY_URL" with
    | null ->
      printfn "SERVICE_DISCOVERY_URL ENV variable not set!"
      exit 1
    | x -> x

  let myIP =
    match Environment.GetEnvironmentVariable "MY_IP" with
    | null ->
      printfn "MY_IP ENV variable not set!"
      exit 1
    | x -> x

  use client = new HttpClient()

  use req =
    new HttpRequestMessage(HttpMethod.Post, "http://" + serviceDiscoveryURL + "/registry/catalog/" + myIP + ":5050")

  try
    client.Send(req) |> ignore
  with _ ->
    ()

  Host
    .CreateDefaultBuilder()
    .ConfigureWebHostDefaults(fun webHost ->
      webHost
        .Configure(configureApp)
        .ConfigureServices(configureServices)
        .UseUrls("http://*:5050")
      |> ignore)
    .Build()
    .Run()

  0
