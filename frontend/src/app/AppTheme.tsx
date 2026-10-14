import { useMemo, type PropsWithChildren } from "react";
import {
  CssBaseline,
  ThemeProvider,
  createTheme,
  useMediaQuery,
} from "@mui/material";

export default function AppTheme({ children }: PropsWithChildren) {
  const prefersDarkMode = useMediaQuery("(prefers-color-scheme: dark)");

  const theme = useMemo(
    () =>
      createTheme({
        palette: {
          mode: prefersDarkMode ? "dark" : "light",
          primary: {
            light: "#b9dcf7",
            main: prefersDarkMode ? "#9acbf2" : "#78afe0",
            dark: "#4f8fc7",
            contrastText: "#102a43",
          },
          background: {
            default: prefersDarkMode ? "#16171d" : "#ffffff",
            paper: prefersDarkMode ? "#1f2028" : "#ffffff",
          },
        },
        typography: {
          fontFamily: 'Roboto, system-ui, "Segoe UI", sans-serif',
        },
      }),
    [prefersDarkMode],
  );

  return (
    <ThemeProvider theme={theme}>
      <CssBaseline />
      {children}
    </ThemeProvider>
  );
}
