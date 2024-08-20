"use client";

import { MainHero } from "@/components/main-hero";
import { ErrorPage } from "@/components/error-page";
import { useEffect, useState } from 'react';
import api from "@/lib/axios";

export default function Home() {
  const [status, setStatus] = useState<string | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const response = await api.get('/status');
        if (response.status === 200) {
          setStatus('ok');
        } else {
          setStatus('error');
        }
      } catch (error) {
        setStatus('error');
      } finally {
        setLoading(false);
      }
    };

    fetchStatus();
  }, []);

  if (loading) {
    return <div>Loading...</div>;
  }

  if (status === 'ok') {
    return (
      <>
        <MainHero />
      </>
    );
  } else {
    return (
      <>
        <ErrorPage />
      </>
    );
  }
}
