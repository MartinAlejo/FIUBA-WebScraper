import { ThemeProvider, createTheme } from "@mui/material";
import Homepage from "./pages/Homepage";


function App() {
  const theme = createTheme({
    palette: {
      mode: "dark",
      primary: {
        main: "#1976d2",
      },
      secondary: {
        main: "#9c27b0",
      },
      error: {
        main: "#d32f2f",
      },
    },
  });

  return (
    <ThemeProvider theme={theme}>
      <Homepage />
    </ThemeProvider>
  );
}

export default App;
