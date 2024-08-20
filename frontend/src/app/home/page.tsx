"use client";

import { useEffect } from "react";
import { useRouter } from "next/navigation";
import { Dashboard } from "@/components/dashboard";
import api from "@/lib/axios";
import Cookies from "js-cookie";

export default function MyDashboard() {
  const router = useRouter();

  useEffect(() => {
    const token = Cookies.get("token");

    if (!token) {
      router.push("/signin"); // Redirect to sign-in page if token is absent
    } else {
      // Fetch user information
      api
        .get("/user/me", {
          headers: {
            Authorization: `${token}`,
          },
        })
        .then((response) => {
          const { username } = response.data;
          sessionStorage.setItem("username", username);
          console.log("User data:", response.data);
        })
        .catch((error) => {
          console.error("Failed to fetch user data:", error);
          router.push("/signin");
        });
    }
  }, [router]);

  return (
    <>
      <Dashboard />
    </>
  );
}
