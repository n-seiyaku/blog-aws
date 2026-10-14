import { Outlet } from "react-router";

export default function AuthLayout() {
  return (
    <div className="flex justify-center items-center h-screen w-screen">
      <Outlet />
    </div>
  );
}
