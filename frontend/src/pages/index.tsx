import React from "react";
import { createRoot } from "react-dom/client";
import { RouterConfig } from "./route";
import Header from "./header";
import { RecoilRoot } from "recoil";

export default function Home() {
	if (typeof window === "object") {
		const rootElement = document.getElementById("__next")!;
		const root = createRoot(rootElement);
		root.render(
			<React.StrictMode>
				<RecoilRoot>
					<Header />
					<div className="h-full flex" id="container">
						<RouterConfig />
					</div>
				</RecoilRoot>
			</React.StrictMode>
		);
	}
}