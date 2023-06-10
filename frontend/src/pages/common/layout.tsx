import { RecoilRoot } from "recoil";
import Header from "./header";
import SideNav1 from "./sideNav1";
import SideNav2 from "./sideNav2";

function Layout({
  children,
}: {
  children: React.ReactNode
}) {
  return (
    <html>
      <head />
      <body>
        <RecoilRoot>
        <div className="h-full">
          <Header />
          <div className="h-full flex">
            <div>
              <SideNav1 />
            </div>
            <div>
              <SideNav2 />
            </div>
            <main>
              { children }
            </main>        
          </div>
          </div>
          </RecoilRoot>
      </body>
    </html>
  )
}

export default Layout;
