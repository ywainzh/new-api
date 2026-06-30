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
import {
  Plus,
  MoreHorizontal,
  RefreshCw,
  List,
  Building2,
  AlertCircle,
} from 'lucide-react'
import { useMutation, useQueryClient } from '@tanstack/react-query'
import { useTranslation } from 'react-i18next'
import { toast } from 'sonner'
import { Button } from '@/components/ui/button'
import {
  DropdownMenu,
  DropdownMenuContent,
  DropdownMenuItem,
  DropdownMenuSeparator,
  DropdownMenuShortcut,
  DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu'
import { syncChannelModels } from '../api'
import { modelsQueryKeys, vendorsQueryKeys } from '../lib'
import { useModels } from './models-provider'

export function ModelsPrimaryButtons() {
  const { t } = useTranslation()
  const queryClient = useQueryClient()
  const { setOpen, setCurrentRow } = useModels()

  const syncChannelModelsMutation = useMutation({
    mutationFn: () => syncChannelModels({ locale: 'zh' }),
    onSuccess: (response) => {
      if (!response.success) {
        toast.error(response.message || t('Operation failed'))
        return
      }

      const createdModels = response.data?.created_models || 0
      const createdBasicModels = response.data?.created_basic_models || 0
      const createdVendors = response.data?.created_vendors || 0
      const skippedModels = response.data?.skipped_models?.length || 0
      const details = [
        t('Created {{count}} model(s)', { count: createdModels }),
        t('Created {{count}} basic model(s)', { count: createdBasicModels }),
        t('Created {{count}} vendor(s)', { count: createdVendors }),
      ]

      if (skippedModels > 0) {
        details.push(t('Skipped {{count}} model(s)', { count: skippedModels }))
      }

      if (createdModels === 0 && createdVendors === 0 && skippedModels === 0) {
        toast.info(t('No new channel models to sync'))
      } else {
        toast.success(t('Channel model sync completed'), {
          description: details.join(' / '),
        })
      }

      queryClient.invalidateQueries({ queryKey: modelsQueryKeys.lists() })
      queryClient.invalidateQueries({ queryKey: modelsQueryKeys.missing() })
      queryClient.invalidateQueries({ queryKey: vendorsQueryKeys.lists() })
    },
    onError: (error) => {
      toast.error((error as Error)?.message || t('Operation failed'))
    },
  })

  const handleCreateModel = () => {
    setCurrentRow(null)
    setOpen('create-model')
  }

  const handleMissingModels = () => {
    setOpen('missing-models')
  }

  const handleSync = () => {
    setOpen('sync-wizard')
  }

  const handleSyncChannelModels = () => {
    syncChannelModelsMutation.mutate()
  }

  const handlePrefillGroups = () => {
    setOpen('prefill-groups')
  }

  const handleManageVendors = () => {
    setOpen('create-vendor') // Will be a separate vendors management dialog
  }

  return (
    <div className='flex items-center gap-2'>
      {/* Create Model */}
      <Button onClick={handleCreateModel} size='sm'>
        <Plus className='h-4 w-4' />
        {t('Add Model')}
      </Button>

      {/* More Actions */}
      <DropdownMenu>
        <DropdownMenuTrigger render={<Button variant='outline' size='sm' />}>
          <MoreHorizontal className='h-4 w-4' />
        </DropdownMenuTrigger>
        <DropdownMenuContent align='end' className='w-56'>
          <DropdownMenuItem onClick={handleMissingModels}>
            {t('Missing Models')}
            <DropdownMenuShortcut>
              <AlertCircle className='h-4 w-4' />
            </DropdownMenuShortcut>
          </DropdownMenuItem>

          <DropdownMenuItem onClick={handleSync}>
            {t('Sync Upstream')}
            <DropdownMenuShortcut>
              <RefreshCw className='h-4 w-4' />
            </DropdownMenuShortcut>
          </DropdownMenuItem>

          <DropdownMenuItem
            onClick={handleSyncChannelModels}
            disabled={syncChannelModelsMutation.isPending}
          >
            {t('Sync Channel Models')}
            <DropdownMenuShortcut>
              <RefreshCw
                className={
                  syncChannelModelsMutation.isPending
                    ? 'h-4 w-4 animate-spin'
                    : 'h-4 w-4'
                }
              />
            </DropdownMenuShortcut>
          </DropdownMenuItem>

          <DropdownMenuSeparator />

          <DropdownMenuItem onClick={handlePrefillGroups}>
            {t('Prefill Groups')}
            <DropdownMenuShortcut>
              <List className='h-4 w-4' />
            </DropdownMenuShortcut>
          </DropdownMenuItem>

          <DropdownMenuItem onClick={handleManageVendors}>
            {t('Manage Vendors')}
            <DropdownMenuShortcut>
              <Building2 className='h-4 w-4' />
            </DropdownMenuShortcut>
          </DropdownMenuItem>
        </DropdownMenuContent>
      </DropdownMenu>
    </div>
  )
}
