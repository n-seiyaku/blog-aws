import { createBrowserRouter } from "react-router";
import RootLayout from "./layout";
import Login from "../features/login/Login";
import Blog from "../features/blog/Blog";

export const router = createBrowserRouter([
  {
    path: "/",
    Component: RootLayout,
    children: [
      {
        path: "login",
        Component: Login,
      },
      {
        path: "blog",
        Component: Blog,
      },
    ],
  },
]);
