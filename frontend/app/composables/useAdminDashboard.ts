import { computed, ref } from 'vue'
import { useMockDataStore } from '~/repositories/mock/MockDataStore'
import type { MarketplaceOrder } from '~/types/order'
import type { User } from '~/types/auth'
import type { Product } from '~/types/product'

export interface MonthlyMetric {
  month: string
  fullName: string
  sales: number
  orders: number
}

export interface ModerationReport {
  id: string
  targetType: 'product' | 'review' | 'shop'
  targetTitle: string
  reason: string
  reporter: string
  date: string
  status: 'pending' | 'resolved' | 'dismissed'
  severity: 'low' | 'medium' | 'high'
}

export interface SystemService {
  name: string
  status: 'healthy' | 'warning' | 'critical'
  uptime: string
  latency: string
  message?: string
}

export const useAdminDashboard = () => {
  const { users, products, orders, shop } = useMockDataStore()

  // Real-time Total Users
  const totalUsers = computed(() => users.value.length)
  const usersGrowth = '↑ 26% this month'

  // Real-time Total Unique Shops
  const totalShops = computed(() => {
    const shopNames = new Set<string>()
    products.value.forEach(p => {
      if (p.shop) shopNames.add(p.shop)
    })
    if (shop.value && shop.value.name) {
      shopNames.add(shop.value.name)
    }
    return Math.max(shopNames.size, 6)
  })
  const shopsGrowth = '↑ 4 new this week'

  // Real-time Total Listed Products
  // Uses active products in store with baseline platform catalogue
  const totalProducts = computed(() => {
    const baseCount = 477 // baseline catalogue
    return baseCount + products.value.length
  })

  // Real-time Total Orders
  const totalOrders = computed(() => orders.value.length)
  const ordersGrowth = '↑ 18% this month'

  // Real-time Gross Revenue / Total Sales
  const totalSalesAmount = computed(() => {
    return orders.value.reduce((sum, order) => sum + (Number(order.total) || 0), 0)
  })

  const formattedTotalSales = computed(() => {
    return `$${totalSalesAmount.value.toLocaleString()}`
  })
  const salesGrowth = '↑ 31% this month'

  // Last 6 Months calculation (Mar, Apr, May, Jun, Jul, Aug)
  const monthlyMetrics = computed<MonthlyMetric[]>(() => {
    const months = [
      { month: 'Mar', fullName: 'March 2026', baseSales: 4200, baseOrders: 85 },
      { month: 'Apr', fullName: 'April 2026', baseSales: 6800, baseOrders: 130 },
      { month: 'May', fullName: 'May 2026', baseSales: 5900, baseOrders: 115 },
      { month: 'Jun', fullName: 'June 2026', baseSales: 8700, baseOrders: 160 },
      { month: 'Jul', fullName: 'July 2026', baseSales: 9400, baseOrders: 185 },
      { month: 'Aug', fullName: 'August 2026', baseSales: 12100, baseOrders: 228 }
    ]

    // Calculate dynamically added order contributions
    const currentOrders = orders.value
    const recentOrdersCount = Math.max(0, currentOrders.length - 8)
    const recentOrdersSales = currentOrders.slice(8).reduce((sum, o) => sum + (o.total || 0), 0)

    return months.map((m, index) => {
      // Add dynamic orders to current month (August)
      const isCurrentMonth = index === months.length - 1
      const extraOrders = isCurrentMonth ? recentOrdersCount : 0
      const extraSales = isCurrentMonth ? recentOrdersSales * 10 : 0

      return {
        month: m.month,
        fullName: m.fullName,
        sales: m.baseSales + extraSales,
        orders: m.baseOrders + extraOrders
      }
    })
  })

  // Moderation Reports (matches badge 3)
  const moderationReports = ref<ModerationReport[]>([
    {
      id: 'REP-101',
      targetType: 'product',
      targetTitle: 'Designer Silk Shirt',
      reason: 'Suspected unauthorized brand reproduction',
      reporter: 'Atelier North',
      date: '2 hours ago',
      status: 'pending',
      severity: 'high'
    },
    {
      id: 'REP-102',
      targetType: 'review',
      targetTitle: 'Review on Clean everyday sneakers',
      reason: 'Contains off-platform external contact info',
      reporter: 'System Bot',
      date: '5 hours ago',
      status: 'pending',
      severity: 'medium'
    },
    {
      id: 'REP-103',
      targetType: 'shop',
      targetTitle: 'QuickSale Store',
      reason: 'Unverified seller payment credentials',
      reporter: 'Risk Engine',
      date: '1 day ago',
      status: 'pending',
      severity: 'low'
    }
  ])

  const pendingReportsCount = computed(() => {
    return moderationReports.value.filter(r => r.status === 'pending').length
  })

  // Infrastructure / System Services (matches badge 1)
  const systemServices = ref<SystemService[]>([
    {
      name: 'API Gateway & Auth',
      status: 'healthy',
      uptime: '99.98%',
      latency: '24ms'
    },
    {
      name: 'PostgreSQL Database',
      status: 'healthy',
      uptime: '99.99%',
      latency: '4ms'
    },
    {
      name: 'Search Indexer & Vector Store',
      status: 'warning',
      uptime: '98.85%',
      latency: '142ms',
      message: 'High memory load (89% utilization)'
    },
    {
      name: 'Storage & Media CDN',
      status: 'healthy',
      uptime: '100%',
      latency: '18ms'
    }
  ])

  const systemAlertsCount = computed(() => {
    return systemServices.value.filter(s => s.status !== 'healthy').length
  })

  const resolveReport = (id: string) => {
    const report = moderationReports.value.find(r => r.id === id)
    if (report) {
      report.status = 'resolved'
    }
  }

  const dismissReport = (id: string) => {
    const report = moderationReports.value.find(r => r.id === id)
    if (report) {
      report.status = 'dismissed'
    }
  }

  return {
    users,
    products,
    orders,
    totalUsers,
    usersGrowth,
    totalShops,
    shopsGrowth,
    totalProducts,
    totalOrders,
    ordersGrowth,
    totalSalesAmount,
    formattedTotalSales,
    salesGrowth,
    monthlyMetrics,
    moderationReports,
    pendingReportsCount,
    systemServices,
    systemAlertsCount,
    resolveReport,
    dismissReport
  }
}
