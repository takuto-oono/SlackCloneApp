import { RecoilRoot } from "recoil";
import Header from "./header";
import SideNav1 from "./sideNav1";
import SideNav2 from "./sideNav2";

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }: LayoutProps)=>{
  return (
    <div className="h-full">
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
    </div>
  )
}

export default Layout;
