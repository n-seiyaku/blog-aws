import { StrictMode } from "react";
import { createRoot } from "react-dom/client";
import { RouterProvider } from "react-router/dom";

import { router } from "./app/router";
import "./index.css";
import { AuthProvider } from "./auth";
import AppTheme from "./app/AppTheme";

createRoot(document.getElementById("root")!).render(
  <StrictMode>
    <AppTheme>
      <AuthProvider>
        <RouterProvider router={router} />
      </AuthProvider>
    </AppTheme>
  </StrictMode>,
);
