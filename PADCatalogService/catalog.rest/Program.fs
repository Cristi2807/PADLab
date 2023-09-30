open Giraffe
open Microsoft.AspNetCore.Builder
open Microsoft.AspNetCore.Hosting
open Microsoft.Extensions.DependencyInjection
open Microsoft.Extensions.Hosting
open catalog.rest

let webApp =
  choose [ ProductApi.productRoutes; RequestErrors.NOT_FOUND "Not found" ]

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
        .UseUrls("http://*:5050")
      |> ignore)
    .Build()
    .Run()

  0
