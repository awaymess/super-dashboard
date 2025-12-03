'use client';

import { useCallback } from 'react';
import { useDispatch, useSelector } from 'react-redux';
import { RootState, AppDispatch } from '@/store';
import { loginStart, loginSuccess, loginFailure, logout, updateUser, clearError } from '@/store/slices/authSlice';
import { User } from '@/types/common';

export function useAuth() {
  const dispatch = useDispatch<AppDispatch>();
  const { user, isAuthenticated, isLoading, error } = useSelector((state: RootState) => state.auth);

  const login = useCallback(
    async (email: string, password: string) => {
      dispatch(loginStart());
      try {
        const mockUser: User = {
          id: '1',
          email,
          name: email.split('@')[0],
          role: 'user',
          preferences: {
            theme: 'dark',
            language: 'en',
            currency: 'USD',
            timezone: 'UTC',
            notifications: {
              email: true,
              push: true,
              valueBets: true,
              priceAlerts: true,
              newsAlerts: false,
              portfolioUpdates: true,
            },
          },
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };
        dispatch(loginSuccess(mockUser));
        return true;
      } catch {
        dispatch(loginFailure('Login failed. Please check your credentials.'));
        return false;
      }
    },
    [dispatch]
  );

  const register = useCallback(
    async (email: string, password: string, name: string) => {
      dispatch(loginStart());
      try {
        const mockUser: User = {
          id: '1',
          email,
          name,
          role: 'user',
          preferences: {
            theme: 'dark',
            language: 'en',
            currency: 'USD',
            timezone: 'UTC',
            notifications: {
              email: true,
              push: true,
              valueBets: true,
              priceAlerts: true,
              newsAlerts: false,
              portfolioUpdates: true,
            },
          },
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
        };
        dispatch(loginSuccess(mockUser));
        return true;
      } catch {
        dispatch(loginFailure('Registration failed. Please try again.'));
        return false;
      }
    },
    [dispatch]
  );

  const logoutUser = useCallback(() => {
    dispatch(logout());
    if (typeof window !== 'undefined') {
      localStorage.removeItem('token');
    }
  }, [dispatch]);

  const updateUserProfile = useCallback(
    (updates: Partial<User>) => {
      dispatch(updateUser(updates));
    },
    [dispatch]
  );

  const resetError = useCallback(() => {
    dispatch(clearError());
  }, [dispatch]);

  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    login,
    register,
    logout: logoutUser,
    updateUser: updateUserProfile,
    clearError: resetError,
  };
}
