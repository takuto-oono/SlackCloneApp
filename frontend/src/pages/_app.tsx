import { AppProps } from "next/app";
import { createContext } from "react";
import { CookiesProvider } from "react-cookie";
import { Link } from "react-router-dom";
import 'src/styles/globals.css'

export default function App({ Component, pageProps }: AppProps) {
  return (
    <CookiesProvider>
      <Component {...pageProps} />
    </CookiesProvider>
  )
}
