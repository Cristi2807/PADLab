namespace catalog

open System.Threading.Tasks

module ServiceUtils =
  open FsToolkit.ErrorHandling
  open Giraffe
  open Microsoft.AspNetCore.Http

  type ServiceError =
    | NotFound of string
    | BadRequest of string
    | InternalError of string

  let inline (=>) a b = a, box b

  let ToAsyncResult<'a> : (Task<'a> -> Async<Result<'a, ServiceError>>) =
    Async.AwaitTask
    >> Async.Catch
    >> Async.map (
      Result.ofChoice
      >> Result.mapError
        // (fun ex -> ex.Message |> InternalError ))
        (fun ex ->
          if ex.InnerException |> isNull then
            ex.ToString() |> InternalError
          else
            "INNER" + ex.InnerException.ToString() + "OUTER" + ex.ToString()
            |> InternalError)
    )

  let UnitTaskToAsyncResult: (Task -> Async<Result<unit, ServiceError>>) =
    Async.AwaitTask
    >> Async.Catch
    >> Async.map (
      Result.ofChoice
      >> Result.mapError
        // (fun ex -> ex.Message |> InternalError ))
        (fun ex ->
          if ex.InnerException |> isNull then
            ex.ToString() |> InternalError
          else
            "INNER" + ex.InnerException.ToString() + "OUTER" + ex.ToString()
            |> InternalError)
    )

  let toApiResponse (ctx: HttpContext) (res: Async<Result<string, ServiceError>>) =
    task {
      let! res' = res

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
    }
