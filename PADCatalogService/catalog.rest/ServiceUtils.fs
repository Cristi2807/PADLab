namespace catalog.rest

module ServiceUtils =
  open FsToolkit.ErrorHandling
  open Giraffe
  open Microsoft.AspNetCore.Http
  open catalog.dal
  open catalog.ServiceUtils

  let toApiResponse (ctx: HttpContext) (dalFactory: DalContextFactory) (res: Async<Result<string, ServiceError>>) =
    task {
      let res' =
        try
          Async.RunSynchronously(res, 15000)
        with _ ->
          dalFactory.Rollback()
          "It takes too long to process the request" |> RequestTimeout |> Error


      ctx.SetHttpHeader("Access-Control-Allow-Methods", "*")
      ctx.SetHttpHeader("Access-Control-Allow-Headers", "*")
      ctx.SetHttpHeader("Access-Control-Allow-Origin", "*")

      match res' with
      | Ok okVal ->
        ctx.SetStatusCode 200
        ctx.SetContentType "application/json"
        return! ctx.WriteStringAsync okVal
      | Error errVal ->
        match errVal with
        | NotFound msg ->
          ctx.SetStatusCode(404)
          return! ctx.WriteJsonAsync msg
        | BadRequest msg ->
          ctx.SetStatusCode(400)
          return! ctx.WriteJsonAsync msg
        | InternalError msg ->
          ctx.SetStatusCode(500)
          return! ctx.WriteJsonAsync msg
        | RequestTimeout msg ->
          ctx.SetStatusCode(408)
          return! ctx.WriteJsonAsync msg
    }
