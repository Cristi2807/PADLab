open Giraffe
open Microsoft.AspNetCore.Builder
open Microsoft.AspNetCore.Hosting
open Microsoft.Extensions.DependencyInjection
open Microsoft.Extensions.Hosting
open inventory.rest
open System
open System.Net.Http

let webApp =
  choose [ InventoryApi.inventoryRoutes; RequestErrors.NOT_FOUND "Not found" ]

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
    new HttpRequestMessage(HttpMethod.Post, "http://" + serviceDiscoveryURL + "/registry/inventory/" + myIP + ":7070")

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
        .UseUrls("http://*:7070")
      |> ignore)
    .Build()
    .Run()

  0
