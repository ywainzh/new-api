/*
Copyright (C) 2023-2026 QuantumNous

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <https://www.gnu.org/licenses/>.

For commercial licensing, please contact support@quantumnous.com
*/
import { createFileRoute, redirect } from '@tanstack/react-router'
import { useEffect } from 'react'

import { Main } from '@/components/layout'
import { Playground } from '@/features/playground'
import { useStatus } from '@/hooks/use-status'
import { getStatus } from '@/lib/api'
import { isSidebarModuleEnabledFromStatus } from '@/lib/nav-modules'

export const Route = createFileRoute('/_authenticated/playground/')({
  beforeLoad: async () => {
    const status = await getStatus()
    if (!isSidebarModuleEnabledFromStatus(status, 'chat', 'playground')) {
      throw redirect({ to: '/dashboard' })
    }
  },
  component: PlaygroundPage,
})

function PlaygroundPage() {
  const { status } = useStatus()

  useEffect(() => {
    if (status?.personal_mode_enabled === true) {
      window.location.replace('/dashboard')
    }
  }, [status?.personal_mode_enabled])

  if (status?.personal_mode_enabled === true) {
    return null
  }

  return (
    <Main className='p-0'>
      <Playground />
    </Main>
  )
}
