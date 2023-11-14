namespace inventory.rest

module ServiceUtils =
  open FsToolkit.ErrorHandling
  open Giraffe
  open Microsoft.AspNetCore.Http
  open inventory.dal
  open inventory.ServiceUtils

  let mutable requestCounts: Map<int, int> =
    Map.empty
    |> Map.add 200 0
    |> Map.add 400 0
    |> Map.add 404 0
    |> Map.add 408 0
    |> Map.add 500 0

  let toApiResponse (ctx: HttpContext) (dalFactory: DalContextFactory) (res: Async<Result<string, ServiceError>>) =
    task {
      let res' =
        try
          Async.RunSynchronously(res, 15000)
        with _ ->
          dalFactory.Rollback()
          dalFactory.Dispose()
          "It takes too long to process the request" |> RequestTimeout |> Error


      ctx.SetHttpHeader("Access-Control-Allow-Methods", "*")
      ctx.SetHttpHeader("Access-Control-Allow-Headers", "*")
      ctx.SetHttpHeader("Access-Control-Allow-Origin", "*")

      match res' with
      | Ok okVal ->
        requestCounts <- requestCounts.Add(200, requestCounts[200] + 1)
        ctx.SetStatusCode 200
        ctx.SetContentType "application/json"
        return! ctx.WriteStringAsync okVal
      | Error errVal ->
        match errVal with
        | NotFound msg ->
          requestCounts <- requestCounts.Add(404, requestCounts[404] + 1)
          ctx.SetStatusCode(404)
          return! ctx.WriteJsonAsync msg
        | BadRequest msg ->
          requestCounts <- requestCounts.Add(400, requestCounts[400] + 1)
          ctx.SetStatusCode(400)
          return! ctx.WriteJsonAsync msg
        | InternalError msg ->
          requestCounts <- requestCounts.Add(500, requestCounts[500] + 1)
          ctx.SetStatusCode(500)
          return! ctx.WriteJsonAsync msg
        | RequestTimeout msg ->
          requestCounts <- requestCounts.Add(408, requestCounts[408] + 1)
          ctx.SetStatusCode(408)
          return! ctx.WriteJsonAsync msg
    }
