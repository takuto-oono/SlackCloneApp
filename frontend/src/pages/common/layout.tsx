import { RecoilRoot } from "recoil";
import Header from "./header";
import SideNav1 from "./sideNav1";
import SideNav2 from "./sideNav2";

interface LayoutProps {
  children: React.ReactNode;
}

function Layout({ children }: LayoutProps){
  return (
    <main className="h-full">
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
            <div>
              { children }
            </div>
          </div>
        </div>
      </RecoilRoot>
    </main>
  )
}

export default Layout;
