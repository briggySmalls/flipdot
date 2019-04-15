import { Machine, StateMachine, interpret } from 'xstate';
import { Interpreter } from 'xstate/lib/interpreter';

export interface AppMachineSchema {
  states: {
    login: {};
    message: {};
    result: {};
  };
}

export type AppMachineEvents =
  | { type: 'AUTH' }
  | { type: 'SENT' }
  | { type: 'REAUTH' }
  | { type: 'NEW' };

// Stateless machine definition
// machine.transition(...) is a pure function used by the interpreter.
export const AppMachine: StateMachine<undefined, AppMachineSchema, AppMachineEvents> = Machine({
  initial: 'login',
  states: {
    login: { on: { AUTH: 'message' } },
    message: { on: { SENT: 'result', REAUTH: 'login' } },
    result: { on: { NEW: 'message' } },
  },
});

export function createFsm(): Interpreter<undefined, AppMachineSchema, AppMachineEvents> {
  return interpret(AppMachine);
}
