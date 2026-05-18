import React from "react";
import ReactDOM from "react-dom/client";
// Suppress TS error for side-effect CSS import when no type declarations are present
// @ts-ignore
import "./index.css";
import { App } from "./App";

ReactDOM.createRoot(document.getElementById("root")!).render(
  <React.StrictMode>
    <App />
  </React.StrictMode>,
);
