import { createBrowserRouter } from "react-router";

import AuthLayout from "./layouts/AuthLayout";
import RootLayout from "./layouts/RootLayout";

import Login from "../features/login/Login";
import Blog from "../features/blog/Blog";

export const router = createBrowserRouter([
  {
    Component: AuthLayout,
    children: [
      {
        path: "/login",
        Component: Login,
      },
    ],
  },
  {
    path: "/",
    Component: RootLayout,
    children: [
      {
        path: "blog",
        Component: Blog,
      },
    ],
  },
]);
