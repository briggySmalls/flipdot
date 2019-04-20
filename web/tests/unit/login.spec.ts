import { shallowMount, Wrapper } from '@vue/test-utils';
import { expect } from 'chai';
import Login from '@/components/Login.vue';
import { Client } from '@/ts/client';
import { Interpreter } from 'xstate/lib/interpreter';
import { mock, instance, when, anything, verify } from 'ts-mockito';
import { CombinedVueInstance, Vue } from 'vue/types/vue';
import { grpc } from '@improbable-eng/grpc-web';
import { StateSchema, EventObject } from 'xstate';
import { createLocalVue } from '@vue/test-utils';
import BootstrapVue from 'bootstrap-vue';

describe('Login.vue', () => {
  let mockedClient: Client;
  let client: Client;

  let mockedFsm: Interpreter<{}, StateSchema, EventObject>;
  let fsm: Interpreter<{}, StateSchema, EventObject>;

  let wrapper: Wrapper<CombinedVueInstance<Login, object, object, object, Record<never, any>>>;

  beforeEach(() => {
    // create an extended `Vue` constructor
    const localVue = createLocalVue();
    // install plugins as normal
    localVue.use(BootstrapVue);
    // Create a mock client
    mockedClient = mock(Client);
    client = instance(mockedClient);
    // Create a mock fsm
    mockedFsm = mock(Interpreter);
    fsm = instance(mockedFsm);
    // Create object under test
    wrapper = shallowMount(Login, {
      localVue,
      propsData: {
        client,
        fsm,
      },
    });
  });

  it('authenticates', () => {
    // Configure client
    when(mockedClient.error).thenReturn(null);
    when(mockedClient.authenticate('password', anything())).thenCall(
        (password: string, callback: (response: any) => void) => {
            // Execute callback
            callback(null);
        },
    );
    // Send password
    wrapper.find('form .login-password').setValue('password');
    wrapper.find('form .login-submit').trigger('submit');
    // Set a password
    verify(mockedClient.authenticate('password', anything())).once();
    verify(mockedFsm.send('AUTH')).once();
  });

  it('displays error', () => {
    // // Configure client
    when(mockedClient.error).thenReturn(null);
    when(mockedClient.authenticate('password', anything())).thenCall(
      (password: string, callback: (response: any) => void) => {
        // Indicate an error occurred
        when(mockedClient.error).thenReturn({
          code: grpc.Code.Unauthenticated,
          message: 'Failed to login',
        });
        // Execute callback
        callback(null);
      },
    );
    // Send password
    wrapper.find('form .login-password').setValue('password');
    wrapper.find('form .login-submit').trigger('submit');
    // Set a password
    verify(mockedClient.authenticate('password', anything())).once();
    verify(mockedFsm.send(anything())).never();
    expect(wrapper.find('.alert').isVisible());
  });
});
