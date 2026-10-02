import { Link, Outlet } from "react-router";
import { useAuth } from "../auth";

export default function RootLayout() {
  const { handleLogout, user } = useAuth();
  return (
    <>
      <header>
        <nav>
          <Link to="/login">Login</Link>
          {" | "}
          <Link to="/blog">Blog</Link>
          {user && (
            <>
              {" | "}
              <button onClick={() => handleLogout()}>Logout</button>
            </>
          )}
        </nav>
      </header>

      <main>
        <Outlet />
      </main>

      <footer>My AWS Blog</footer>
    </>
  );
}
