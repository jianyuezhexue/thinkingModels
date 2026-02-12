<script lang="ts" setup>
import type { VbenFormSchema } from '@vben/common-ui';
import type { Recordable } from '@vben/types';

import { computed, h, ref } from 'vue';
import { useRouter } from 'vue-router';

import { AuthenticationRegister, z } from '@vben/common-ui';
import { $t } from '@vben/locales';

import { registerApi } from '#/api/core/auth';
import { ElMessage } from 'element-plus';

defineOptions({ name: 'Register' });

const loading = ref(false);
const router = useRouter();

const formSchema = computed((): VbenFormSchema[] => {
  return [
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: '请输入用户名',
      },
      fieldName: 'username',
      label: '用户名',
      rules: z.string().min(3, { message: '用户名至少3个字符' }).max(20, { message: '用户名最多20个字符' }),
    },
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: '请输入昵称（可选）',
      },
      fieldName: 'nickname',
      label: '昵称',
      rules: z.string().max(50, { message: '昵称最多50个字符' }).optional(),
    },
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: '请输入邮箱（可选）',
      },
      fieldName: 'email',
      label: '邮箱',
      rules: z.string().email({ message: '请输入有效的邮箱地址' }).optional().or(z.literal('')),
    },
    {
      component: 'VbenInput',
      componentProps: {
        placeholder: '请输入手机号（可选）',
      },
      fieldName: 'phone',
      label: '手机号',
      rules: z.string().length(11, { message: '手机号必须为11位' }).regex(/^1[3-9]\d{9}$/, { message: '请输入有效的手机号' }).optional().or(z.literal('')),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        passwordStrength: true,
        placeholder: '请输入密码',
      },
      fieldName: 'password',
      label: '密码',
      renderComponentContent() {
        return {
          strengthText: () => $t('authentication.passwordStrength'),
        };
      },
      rules: z.string().min(8, { message: '密码至少8个字符' }),
    },
    {
      component: 'VbenInputPassword',
      componentProps: {
        placeholder: '请再次输入密码',
      },
      dependencies: {
        rules(values) {
          const { password } = values;
          return z
            .string({ required_error: '请确认密码' })
            .min(1, { message: '请确认密码' })
            .refine((value) => value === password, {
              message: '两次输入的密码不一致',
            });
        },
        triggerFields: ['password'],
      },
      fieldName: 'confirmPassword',
      label: '确认密码',
    },
    {
      component: 'VbenCheckbox',
      fieldName: 'agreePolicy',
      renderComponentContent: () => ({
        default: () =>
          h('span', [
            $t('authentication.agree'),
            h(
              'a',
              {
                class: 'vben-link ml-1 ',
                href: '',
              },
              `${$t('authentication.privacyPolicy')} & ${$t('authentication.terms')}`,
            ),
          ]),
      }),
      rules: z.boolean().refine((value) => !!value, {
        message: $t('authentication.agreeTip'),
      }),
    },
  ];
});

async function handleSubmit(values: Recordable<any>) {
  loading.value = true;
  try {
    // 构造注册参数
    const registerParams = {
      username: values.username,
      password: values.password,
      nickname: values.nickname || '',
      email: values.email || '',
      phone: values.phone || '',
    };

    // 调用注册API
    const result = await registerApi(registerParams);

    ElMessage.success('注册成功，请登录');

    // 注册成功后跳转到登录页面
    setTimeout(() => {
      router.push('/auth/login');
    }, 1500);

    // eslint-disable-next-line no-console
    console.log('注册成功:', result);
  } catch (error: any) {
    ElMessage.error(error?.message || '注册失败，请重试');
  } finally {
    loading.value = false;
  }
}
</script>

<template>
  <AuthenticationRegister
    :form-schema="formSchema"
    :loading="loading"
    @submit="handleSubmit"
  />
</template>
