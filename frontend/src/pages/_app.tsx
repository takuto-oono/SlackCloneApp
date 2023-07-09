import { AppProps } from "next/app";
import { CookiesProvider } from "react-cookie";
import 'src/styles/globals.css'
import Layout from "./common/layout";
import { RecoilRoot } from "recoil";

export default function App({ Component, pageProps }: AppProps) {
  return (
    <CookiesProvider>
      <RecoilRoot>
        <Layout>
          <Component {...pageProps} />
        </Layout>
      </RecoilRoot>
    </CookiesProvider>
  )
}
