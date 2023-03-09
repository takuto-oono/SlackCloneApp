
import { useRouter } from "next/router";

function ShowWorkspace() {
    const router = useRouter();

  return (
    <div>
      {" " + router.query.id}
    </div>
  )
}

export default ShowWorkspace;
