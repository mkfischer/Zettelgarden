import React, { useState, useEffect } from "react";
import "../App.css";
import { SearchPage } from "./cards/SearchPage";
import { UserSettingsPage } from "./UserSettings";
import { FileVault } from "./FileVault";
import { ViewPage } from "./cards/ViewPage";
import { EditPage } from "./cards/EditPage";
import { Sidebar } from "../components/Sidebar";
import { useAuth } from "../contexts/AuthContext";
import { Navigate, useNavigate } from "react-router-dom";
import { Route, Routes } from "react-router-dom";
import { EmailValidationBanner } from "../components/EmailValidationBanner";
import { BillingSuccess } from "./BillingSuccess";
import { BillingCancelled } from "./BillingCancelled";
import { SubscriptionPage } from "./SubscriptionPage";
import { DashboardPage } from "./DashboardPage";

import { Card } from "../models/Card";
import { TaskPage } from "./tasks/TaskPage";
import { TaskProvider, useTaskContext } from "../contexts/TaskContext";
import {
  PartialCardProvider,
  usePartialCardContext,
} from "../contexts/CardContext";

function MainAppContent() {
  const navigate = useNavigate();
  const [lastCardId, setLastCardId] = useState("");
  const [searchTerm, setSearchTerm] = useState("");
  const [searchCards, setSearchCards] = useState<Card[]>([]);
  const { isAuthenticated, isLoading, hasSubscription, logoutUser } = useAuth();
  const { setRefreshTasks } = useTaskContext();
  const { setRefreshPartialCards } = usePartialCardContext();

  // changing pages

  async function handleNewCard(cardType: string) {
    navigate("/app/card/new", { state: { cardType: cardType } });
  }

  useEffect(() => {
    if (!localStorage.getItem("token")) {
      logoutUser();
    }
  }, [isAuthenticated]);

  useEffect(() => {
    setRefreshTasks(true);
    setRefreshPartialCards(true);
  }, []);

  if (!isLoading) {
    if (!isAuthenticated) {
      <Navigate to="/login" />;
    }
    return (
      <div className="main-content">
        <Sidebar />
        <div className="content">
          <div className="content-display">
            <EmailValidationBanner />
            <Routes>
              {!hasSubscription && (
                <>
                  <Route path="subscription" element={<SubscriptionPage />} />
                  <Route
                    path="settings/billing/success"
                    element={<BillingSuccess />}
                  />
                  <Route
                    path="settings/billing/cancelled"
                    element={<BillingCancelled />}
                  />
                </>
              )}
              {hasSubscription ? (
                <>
                  <Route
                    path="search"
                    element={
                      <SearchPage
                        searchTerm={searchTerm}
                        setSearchTerm={setSearchTerm}
                        cards={searchCards}
                        setCards={setSearchCards}
                      />
                    }
                  />
                  <Route
                    path="card/:id"
                    element={<ViewPage setLastCardId={setLastCardId} />}
                  />
                  <Route
                    path="card/:id/edit"
                    element={
                      <EditPage newCard={false} lastCardId={lastCardId} />
                    }
                  />

                  <Route
                    path="card/new"
                    element={
                      <EditPage newCard={true} lastCardId={lastCardId} />
                    }
                  />
                  <Route path="settings" element={<UserSettingsPage />} />
                  <Route path="files" element={<FileVault />} />
                  <Route path="tasks" element={<TaskPage />} />
                  <Route path="*" element={<DashboardPage />} />
                </>
              ) : (
                <Route
                  path="*"
                  element={<Navigate to="/app/subscription" replace />}
                />
              )}
            </Routes>
          </div>
        </div>
      </div>
    );
  }
}

function MainApp() {
  return (
    <PartialCardProvider>
      <TaskProvider>
        <MainAppContent />
      </TaskProvider>
    </PartialCardProvider>
  );
}

export default MainApp;
