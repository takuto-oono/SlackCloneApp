import { AppProps } from "next/app";
import { ReactElement, ReactNode } from "react";
import { CookiesProvider } from "react-cookie";
import 'src/styles/globals.css'
import Layout from "./common/layout";
import { NextPage } from "next/types";
import { RecoilRoot } from "recoil";

 export type NextPageWithLayout<P = {}, IP = P> = NextPage<P, IP> & {
    getLayout?: (page: ReactElement) => ReactNode
  }
  
  type AppPropsWithLayout = AppProps & {
    Component: NextPageWithLayout
  }

export default function App({ Component, pageProps }: AppPropsWithLayout) {
  return (
    <CookiesProvider>
      <Component {...pageProps} />
    </CookiesProvider>
  )
}
