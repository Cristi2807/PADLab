open Giraffe
open Microsoft.AspNetCore.Builder
open Microsoft.AspNetCore.Hosting
open Microsoft.Extensions.DependencyInjection
open Microsoft.Extensions.Hosting
open inventory.rest

let webApp =
  choose [ InventoryApi.inventoryRoutes; RequestErrors.NOT_FOUND "Not found" ]

let configureApp (app: IApplicationBuilder) = app.UseGiraffe webApp

let configureServices (services: IServiceCollection) = services.AddGiraffe() |> ignore

[<EntryPoint>]
let main _ =
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
