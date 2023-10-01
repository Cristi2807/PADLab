namespace catalog

open System.Threading.Tasks

module ServiceUtils =
  open FsToolkit.ErrorHandling
  open Microsoft.AspNetCore.Http

  type ServiceError =
    | NotFound of string
    | BadRequest of string
    | InternalError of string
    | RequestTimeout of string

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
