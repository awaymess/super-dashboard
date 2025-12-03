'use client';

import { Target, CheckCircle } from 'lucide-react';
import { GlassCard, Progress, Badge } from '@/components/ui';

interface Goal {
  id: string;
  title: string;
  target: number;
  current: number;
  unit: string;
  deadline?: string;
}

interface GoalTrackerProps {
  goals: Goal[];
}

export function GoalTracker({ goals }: GoalTrackerProps) {
  return (
    <GlassCard>
      <div className="flex items-center gap-2 mb-4">
        <Target className="w-5 h-5 text-primary" />
        <h3 className="font-bold text-white">Goals</h3>
      </div>

      <div className="space-y-4">
        {goals.map((goal) => {
          const progress = Math.min((goal.current / goal.target) * 100, 100);
          const isCompleted = progress >= 100;

          return (
            <div key={goal.id} className="p-4 bg-white/5 rounded-lg">
              <div className="flex items-start justify-between mb-2">
                <div className="flex items-center gap-2">
                  {isCompleted && <CheckCircle className="w-4 h-4 text-success" />}
                  <span className={`font-medium ${isCompleted ? 'text-success' : 'text-white'}`}>
                    {goal.title}
                  </span>
                </div>
                {goal.deadline && (
                  <Badge variant={isCompleted ? 'success' : 'default'}>{goal.deadline}</Badge>
                )}
              </div>
              
              <div className="mb-2">
                <Progress 
                  value={goal.current} 
                  max={goal.target} 
                  color={isCompleted ? 'success' : 'primary'} 
                />
              </div>
              
              <div className="flex justify-between text-sm">
                <span className="text-gray-400">
                  {goal.current.toLocaleString()} / {goal.target.toLocaleString()} {goal.unit}
                </span>
                <span className={isCompleted ? 'text-success font-medium' : 'text-gray-400'}>
                  {progress.toFixed(0)}%
                </span>
              </div>
            </div>
          );
        })}
      </div>
    </GlassCard>
  );
}
