import { friendlyError } from "./errors";

test("известная ошибка переводится в человеческое сообщение", () => {
  const err = new Error("permission denied");
  expect(friendlyError(err)).toBe("Недостаточно прав");
});

test("неизвестная ошибка падает в fallback", () => {
  const err = new Error("kaboom");
  expect(friendlyError(err, "что-то пошло не так")).toBe("что-то пошло не так");
});

test("если пришёл не Error — берётся fallback", () => {
  expect(friendlyError("plain string", "fallback")).toBe("fallback");
});
