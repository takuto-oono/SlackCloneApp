import React from "react";
import { createRoot } from "react-dom/client";
import { RecoilRoot } from "recoil";
import Main from "./main";
import Layout from "./common/layout";

export default function Home() {
	if (typeof window === "object") {
		const rootElement = document.getElementById("__next")!;
		const root = createRoot(rootElement);
		root.render(
			<RecoilRoot>
        <Layout>
          <Main />
        </Layout>
      </RecoilRoot>
		);
	}
}
