import { shallowMount } from '@vue/test-utils';
import Login from '@/components/Login.vue';
import { Client } from '@/ts/client';
import { Interpreter } from 'xstate/lib/interpreter';
import { mock, instance, when, anything, verify } from 'ts-mockito';

const mockedFsm = mock(Interpreter);

describe('Login.vue', () => {
  it('authenticates', () => {
    // Create a mock client
    const mockedClient = mock(Client);
    const client = instance(mockedClient);
    when(mockedClient.error).thenReturn(null);
    when(mockedClient.authenticate('password', anything())).thenCall(
        (password: string, callback: (response: any) => void) => {
            // Execute callback
            callback(null);
        },
    );
    // Create a mock fsm
    const fsm = instance(mockedFsm);
    // Create the component
    const wrapper = shallowMount(Login, {
        propsData: {
            client,
            fsm,
        },
    });
    // Send password
    wrapper.find('#login-password').setValue('password');
    wrapper.find('#login-submit').trigger('submit');
    // Set a password
    verify(mockedClient.authenticate('password', anything())).once();
    verify(mockedFsm.send('AUTH')).once();
  });
});
